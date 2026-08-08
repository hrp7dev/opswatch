package ui
import(
	"fmt"
	"strings"
)
var cpuHistory []float64
func AddCPUHistory(value float64){
	if len(cpuHistory)>=50{
		cpuHistory=cpuHistory[1:]
	}
	cpuHistory=append(cpuHistory,value)
}
func RenderCPUChart()string{
	width:=50
	var rows []string
	history:=cpuHistory
	if len(history)>width{
		history=history[len(history)-width:]
	}
	for level:=100;level>=0;level-=10{
		var row strings.Builder
		row.WriteString(fmt.Sprintf("%3d%% │ ",level))
		for i:=0;i<width;i++{
			if i<len(history){
				value:=history[i]
				if value>=float64(level){
					row.WriteString("█ ")
				}else{
					row.WriteString("  ")
				}
			}else{
				row.WriteString("  ")
			}
		}
		rows=append(rows,row.String())
	}
	rows=append(rows,"     └"+strings.Repeat("──",width))
	rows=append(rows,"      "+strings.Repeat(" ",width)+"TIME →")
	return strings.Join(rows,"\n")
}