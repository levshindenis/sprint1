package storages

import (
	"github.com/levshindenis/sprint1/internal/app/config"
	"github.com/levshindenis/sprint1/internal/app/tools"
)

type ServerStorage struct {
	st Storage
	sa config.ServerAddress
}

func (serv *ServerStorage) Init() {
	serv.st.EmptyStorage()
	serv.sa.ParseFlags()
}

func (serv *ServerStorage) InitStorage() {
	serv.st.EmptyStorage()
}

func (serv *ServerStorage) GetStartSA() string {
	return serv.sa.GetStartAddress()
}

func (serv *ServerStorage) GetBaseSA() string {
	return serv.sa.GetShortBaseURL()
}

func (serv *ServerStorage) GetStorage() Storage {
	return serv.st
}

func (serv *ServerStorage) SetStorage(key string, value string) {
	serv.st[key] = value
}

func (serv *ServerStorage) SetBaseSA(value string) {
	serv.sa.SetShortBaseURL(value)
}

func (serv *ServerStorage) ValueInStorage(str string) (string, bool) {
	return serv.st.ValueIn(str)
}

// GetAddress проверяю, есть ли такой адрес(длинный URL) в storage. Если есть, то возвращаю уже заданный короткий,
// если нет, то создаю короткий URL и возвращаю его
func (serv *ServerStorage) GetAddress(str string) string {
	addr := serv.GetBaseSA() + "/"

	if value, ok := serv.ValueInStorage(str); ok {
		return addr + value
	} else {
		shortKey := tools.GenerateShortKey()
		for {
			if _, in := (serv.GetStorage())[shortKey]; !in {
				serv.SetStorage(shortKey, str)
				break
			}
			shortKey = tools.GenerateShortKey()
		}
		return addr + shortKey
	}
}
