package kmgTime

import "time"

var BeijingZone = time.FixedZone("CST", 8*60*60)
var ESTZone = time.FixedZone("EST", -5*60*60) // Eastern Standard Time(加拿大)

//golang 的时区实现看上去很复杂,而且有系统依赖,此处添加一个简单时区枚举,不考虑夏令时,也不考虑时区变化
func MustLoadZone(name string) (loc *time.Location) {
	switch name {
	case "CST", "Beijing":
		return BeijingZone
	case "EST":
		return ESTZone
	default:
		panic("time zone name [" + name + "]not found")
	}
	return nil
}
