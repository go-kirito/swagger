package swagger

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/go-kirito/pkg/api/metadata"
	"google.golang.org/grpc"

	"github.com/go-kirito/pkg/application"

	_ "financial/third_party/swagger/swagger_ui/statik" // import statik static files

	"github.com/rakyll/statik/fs"
)

func Start(app *application.App) error {

	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	sh := http.StripPrefix("/q/swagger-ui", staticServer)

	app.HttpServer().HandlePrefix("/q/swagger-ui", sh)
	app.HttpServer().HandlePrefix("/q/swagger", http.HandlerFunc(swaggerFile))
	app.HttpServer().Handle("/q/services", GetServices(nil))

	return nil
}

func swaggerFile(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "swagger.json") {
		log.Printf("Not Found: %s", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	p := strings.TrimPrefix(r.URL.Path, "/q/swagger/")
	name := path.Join(".", p)
	http.ServeFile(w, r, name)
}

func GetServices(srv *grpc.Server) http.HandlerFunc {
	s := metadata.NewServer(srv)
	return func(w http.ResponseWriter, r *http.Request) {
		reply, err := s.ListServices(context.Background(), &metadata.ListServicesRequest{})
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		var files []string
		for _, service := range reply.GetServices() {
			if service == "kratos.api.Metadata" ||
				service == "grpc.reflection.v1alpha.ServerReflection" ||
				service == "grpc.health.v1.Health" {
				continue
			}

			pathList := strings.Split(service, ".")
			path := strings.ToLower(strings.Join(pathList[:len(pathList)], "/")) + ".swagger.json"
			files = append(files, path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(files)
	}
}
