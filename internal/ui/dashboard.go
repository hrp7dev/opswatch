package ui
import(
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
	"github.com/hrp7dev/opswatch/internal/system"
)
var(
	titleStyle=lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#38BDF8"))
	subtitleStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	labelStyle=lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#CBD5E1"))
	valueStyle=lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F8FAFC"))
	infoStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#94A3B8"))
	panelStyle=lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#334155")).Padding(1,2).Width(32)
	chartStyle=lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#475569")).Padding(1,2).Width(110)
	successStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#22C55E"))
	warningStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#EAB308"))
	dangerStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444"))
	emptyBarStyle=lipgloss.NewStyle().Foreground(lipgloss.Color("#1E293B"))
)

func RenderDashboard(cpuUsage float64,memory system.MemoryStats,disk system.DiskStats){
	Clear()

	header:=lipgloss.JoinVertical(
		lipgloss.Center,
		titleStyle.Render("⚡ OPSWATCH"),
		subtitleStyle.Render("Real-time Server Monitoring"),
		infoStyle.Render("CPU • MEMORY • DISK • LIVE METRICS"),
	)

	cpu:=panelStyle.Render(renderMetric(
		"CPU",
		cpuUsage,
		fmt.Sprintf("%d Cores",system.CPUCores()),
	))

	memoryMetric:=panelStyle.Render(renderMetric(
		"MEMORY",
		memory.Usage,
		fmt.Sprintf("%s / %s",system.FormatBytes(memory.Used),system.FormatBytes(memory.Total)),
	))

	diskMetric:=panelStyle.Render(renderMetric(
		"DISK",
		disk.Usage,
		fmt.Sprintf("%s / %s",system.FormatBytes(disk.Used),system.FormatBytes(disk.Total)),
	))

	stats:=lipgloss.JoinHorizontal(
		lipgloss.Top,
		cpu,
		"  ",
		memoryMetric,
		"  ",
		diskMetric,
	)

	chart:=chartStyle.Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			labelStyle.Render("CPU HISTORY"),
			"",
			RenderCPUChart(),
		),
	)

	body:=lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"",
		stats,
		"",
		chart,
	)

	fmt.Println(body)
}

func renderMetric(label string,percentage float64,detail string)string{
	bar:=progressBar(percentage,20)
	value:=valueStyle.Render(fmt.Sprintf("%6.2f%%",percentage))
	return lipgloss.JoinVertical(
		lipgloss.Left,
		labelStyle.Render(label),
		fmt.Sprintf("%s %s",bar,value),
		infoStyle.Render(detail),
	)
}

func progressBar(percentage float64,width int)string{
	if percentage<0{
		percentage=0
	}
	if percentage>100{
		percentage=100
	}
	filled:=int((percentage/100)*float64(width))
	style:=successStyle
	if percentage>=70{
		style=warningStyle
	}
	if percentage>=90{
		style=dangerStyle
	}
	return style.Render(strings.Repeat("█",filled))+emptyBarStyle.Render(strings.Repeat("░",width-filled))
}