#Базовый сервис для мониторинга докер контейнеров#
![image](https://github.com/user-attachments/assets/f1763fba-d5e0-41e6-9d29-97d9f4706812)

**Команды для запуска**
`docker compose build`
`docker compose up`

**FRONTEND**
`localhost:3000`

**Endpoints**
`r.GET("/containers", handler.GetAllContainers)`
`r.GET("/container/:id", handler.GetContainerByID)`
`r.POST("/container", handler.CreateContainer)`
`r.PUT("/container", handler.UpdateContainerByID)`
`r.DELETE("/container/:id", handler.DeleteContainerByID)`
`r.PUT("/pinger", handler.UpdateTimeContainers) // Endpoint for update time pings for pinger`

 **models.Container**
 `type Container struct {`
`	Id              int       `json:"id"``
`	Address         string    `json:"address"``
`	LastSuccessPing time.Time `json:"last_success_ping"``
`	LastPing        time.Time `json:"last_ping"``
`}`

![image](https://github.com/user-attachments/assets/a1faf6c8-f063-48b2-a6ec-b492384cd434)

###Реализация Pinger###
Сервис для пинга IP адресов основан на паттерне Worker Pool. Worker'ы получают через канал слайс контейнеров, пингуют каждый из них. Пинги сделал через установку tcp соединения, если получилось установить, то контейнер работает.

###REST API###
**Framework** Gin
**Database** PostgreSQL

**Endpoint UpdateTimeContainers отправляет Batch в PostgreSQL**
Pinger отправляет на запись большое количество данных, поэтому чтобы не подключаться много раз каждый воркер после завершения отправляет слайс "опрошенных" контенеров в JSON. REST API используя pgx.Batch отправляет записывает сразу набор данных.



