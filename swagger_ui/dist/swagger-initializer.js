window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">
  const servicesUrl = new URL("/q/services", window.location.href);
  fetch(servicesUrl.toString())
      .then(response => response.json())
      .then(data => {
        const urls = data.map((x) => {
          const url = new URL("/q/swagger/" + x, window.location.href);
          return {url: url.toString(), name: x}
        });
        // the following lines will be replaced by docker/configurator, when it runs in a docker-container
        window.ui = SwaggerUIBundle({
          urls: urls,
          dom_id: '#swagger-ui',
          deepLinking: true,
          presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
          ],
          plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
          ],
          layout: "StandaloneLayout"
        });
      });
  //</editor-fold>
};
