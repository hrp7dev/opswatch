package main
import(
	"fmt"
	"os"
	"log"
	"time"
	"github.com/hrp7dev/opswatch/internal/version"
	"github.com/hrp7dev/opswatch/internal/system"
	"github.com/hrp7dev/opswatch/internal/ui"
)
func main(){
	if len(os.Args)>1&&os.Args[1]=="--version"{
		fmt.Println("OpsWatch",version.Version)
		return
	}
	
	for{
		cpuUsage,err:=system.CPUUsage()
		if err!=nil{
			log.Println("CPU:",err)
			continue
		}
		memory,err:=system.MemoryUsage()
		if err!=nil{
			log.Println("MEMORY:",err)
			continue
		}
		disk,err:=system.DiskUsage("/")
		if err!=nil{
			log.Println("DISK:",err)
			continue
		}
		ui.AddCPUHistory(cpuUsage)
		ui.Clear()
		ui.RenderDashboard(cpuUsage,memory,disk)
		time.Sleep(time.Second)
	}
}