package prefabrely

import (
	"fmt"

	"github.com/372572571/Exercise/utils/compip"
)

// Prefab ...
type Prefab struct{}

// Handle 处理请求
func (h Prefab) Handle(msg *compip.OnMsg, service *compip.Service) {
	fmt.Println(msg.Data)
}
