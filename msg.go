package main

const (
	errDecodeJson = "Не смогли декодировать json: %v\n"
	errSaveItem   = "Не смогли сохранить информацию\n"
	errMethod     = "Не соответствие метода %s\n"
)

const (
	msgNewRequest    = "Получен запрос с url: %s, методом: %s, по протоколу %s\n"
	msgSaveSucsecc   = "Успешное сохранение в БД %s\n"
	msgUpdateSuccess = "Успешное обновление в БД %s\n"
	msgSelectSuccess = "Успешное отображение %s\n"
	msgSaveSuccess   = "Успешная запись в БД"
)

const (
	dir  = "./data"
	perm = 0644
)
