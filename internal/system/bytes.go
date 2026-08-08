package system
import "fmt"
func FormatBytes(bytes uint64)string{
	const unit=1024
	if bytes<unit{
		return fmt.Sprintf("%d B",bytes)
	}
	value:=float64(bytes)
	units:=[]string{"B","KB","MB","GB","TB","PB"}
	i:=0
	for value>=unit&&i<len(units)-1{
		value/=unit
		i++
	}
	return fmt.Sprintf("%.2f %s",value,units[i])
}