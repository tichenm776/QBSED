package service

//
import (
	"context"
	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"zhiyuan/QBSED/configs"
	//"zhiyuan/ai_dormitory_apis/school_affairs/conf"
	"zhiyuan/QBSED/internal/dao"
)

//// Service service.
type Service struct {
	//	ac  *paladin.Map
	//	dao *dao.Dao
	c   *configs.Config
	dao *dao.Dao
}

//
// New new a service and return.
func New(c *configs.Config) (s *Service) {
	//var ac = new(paladin.TOML)
	//if err := paladin.Watch("application.toml", ac); err != nil {
	//	panic(err)
	//}
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return s
}

//
// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) (err error) {
	return s.dao.Ping(ctx)
}

//
// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) ReadPorttoml() (int, bool) {

	type Port_obj struct {
		Port int
	}
	var port Port_obj
	path_GI := "./port.toml"

	if _, err := toml.DecodeFile(path_GI, &port); err != nil {
		log.Error("read toml file error ", err)
		return 0, false
	}

	return port.Port, true
}

//func (s *Service) Saveconfigtoml(setconfig model.Set_config) (bool) {
//
//	file_path := "./setconfig.toml"
//	f, err := os.Create(file_path)
//	if err != nil {
//		log.Fatal(err)
//		return false
//	}
//	if err := toml.NewEncoder(f).Encode(setconfig); err != nil {
//		// failed to encode
//		log.Fatal(err)
//		return false
//	}
//	if err := f.Close(); err != nil {
//		// failed to close the file
//		log.Fatal(err)
//		return false
//	}
//	return true
//}
//
//func (s *Service) Readconfigtoml() (config model.Set_config) {
//
//	file_path := "./setconfig.toml"
//	if _, err := toml.DecodeFile(file_path, &config); err != nil {
//		log.Error("read toml file error(%v)", err)
//	}
//	return config
//}

//
//func (s *Service) GeneratePort(address string) (string) {
//	starter := address
//	numberarr := strings.Split(starter, ".")
//	variable, _ := strconv.Atoi(numberarr[3])
//	//sort.Strings(numberarr)
//	//larger,_ := strconv.Atoi(numberarr[3])
//	PORT := variable * variable
//	if PORT < 10000 {
//		PORT += 10000
//	}
//	return strconv.Itoa(PORT)
//}
//
//func (s *Service) FindUrl(str string) string {
//	// 创建一个正则表达式匹配规则对象
//	if str == "" {
//		return ""
//	}
//	regular := `(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)\.(25[0-5]|2[0-4]\d|[0-1]\d{2}|[1-9]?\d)`
//	reg := regexp.MustCompile(regular)
//	// 利用正则表达式匹配规则对象匹配指定字符串
//	res := reg.FindAllString(str, -1)
//	if (res == nil) {
//		return "3.14"
//	}
//	ip := strings.Join(res, ".")
//	return ip
//}
//
//func (s *Service) Killserver(pid string)(bool) {
//
//	command := "/home/zybox/device_manager/stop.sh"
//	cmd := exec.Command(command, pid)
//	stdout, err := cmd.Output()
//	cmd.Run()
//	log.Println("kill client server")
//	if err != nil {
//		fmt.Println(err.Error())
//		return false
//	}
//	fmt.Print(string(stdout))
//	return true
//}
