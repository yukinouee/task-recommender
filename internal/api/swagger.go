package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// NewSwaggerHandler Swagger UIを提供するハンドラーを作成
func NewSwaggerHandler() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("doc.json"), // Swaggerドキュメントのパス
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)
}
