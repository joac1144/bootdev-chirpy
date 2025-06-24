package api

import (
	"fmt"
	"net/http"
)

const MetricsPath string = "GET /admin/metrics"

func (config *ApiConfig) CountHitsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html")
	hits := config.FileserverHits.Load()
	html := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`
	rw.Write([]byte(fmt.Sprintf(html, hits)))
}
