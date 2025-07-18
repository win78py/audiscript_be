package routes_test

// import (
//     // "audiscript_be/internal/handler"
//     "github.com/gin-gonic/gin"
//     "net/http"
//     "net/http/httptest"
//     "testing"
// )

// func TestHelloWorldHandler(t *testing.T) {
//     // Tạo mock handler, có thể truyền nil nếu không cần App
//     h := &handler.Handler{App: nil}
//     r := gin.New()
//     r.GET("/", h.HelloWorld)

//     req, err := http.NewRequest("GET", "/", nil)
//     if err != nil {
//         t.Fatal(err)
//     }
//     rr := httptest.NewRecorder()
//     r.ServeHTTP(rr, req)

//     if status := rr.Code; status != http.StatusOK {
//         t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
//     }
//     expected := "{\"message\":\"Hello World\"}"
//     if rr.Body.String() != expected {
//         t.Errorf("Handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
//     }
// }