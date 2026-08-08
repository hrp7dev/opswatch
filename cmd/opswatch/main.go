package main
import(
	"fmt"
	"os"
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
			ui.AddLog("error", fmt.Sprintf("CPU check failed: %v", err))
			continue
		}
		memory,err:=system.MemoryUsage()
		if err!=nil{
			ui.AddLog("warn", fmt.Sprintf("Memory check failed: %v", err))
			continue
		}
		disk,err:=system.DiskUsage("/")
		if err!=nil{
			ui.AddLog("warn", fmt.Sprintf("Disk check failed: %v", err))
			continue
		}
		ui.AddCPUHistory(cpuUsage)
		ui.Clear()
		ui.RenderDashboard(cpuUsage,memory,disk)
		time.Sleep(time.Second)
	}
}