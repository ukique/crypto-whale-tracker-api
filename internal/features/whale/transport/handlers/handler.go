package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WhaleHandler(w http.ResponseWriter, r *http.Request) {
	r := gin.Default()
	r.Run(":8081")
}
