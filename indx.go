/*
	indx package.
	Пакет реализует API биржи Indx.

*/
package indx

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// URL запроса
const URL string = "https://api.indx.ru/api/v2/trade/"

// Коды выполнения запроса
var IndxCodes map[string]string = map[string]string{
	// код воврата -> описание
	"0":   "Запрос выполнен успешно",
	"-1":  "Сервис остановлен",
	"-2":  "Доступ запрещен",
	"-3":  "Ошибочный WMID Трейдера",
	"-4":  "Подпись запроса сформирована не верно",
	"-5":  "Не корректная",
	"-6":  "Не существующий номер инструмента",
	"-7":  "Вызов веб сервиса завершилось ошибкой",
	"-8":  "Внутренняя ошибка",
	"-9":  "Неизвестная ошибка",
	"-10": "Неизвестная ошибка",
}

type Indx struct {
	Login    string
	Password string
	Wmid     string
	Culture  string
	sign     string
}

type IndxRequest struct {
	method string
	data   string
}

type IndxBalance struct {
	Code  int
	Desc  string
	Value struct {
		Wmid     string
		Nickname string
		Balance  struct {
			Price float32
			Wmz   float32
		}
		Portfolio []struct {
			ID    int
			Name  string
			Notes int
			Price float32
			Type  string
			Kind  int
			By    int
		}
		Profit []struct {
			Symbolid int
			Buy      float32
			Sell     float32
		}
	}
}

type IndxTools struct {
	Code  int
	Desc  string
	Value []struct {
		ID    int
		Name  string
		Price float32
		Kind  int
		Type  string
		By    int
	}
}

type IndxHistoryTrading struct {
	Code  int
	Desc  string
	Value []struct {
		ID    int
		Stamp int
		Name  string
		Isbid int
		Notes int
		Price float32
	}
}

type IndxHistoryTransaction struct {
	Code  int
	Desc  string
	Value []struct {
		Stamp    int
		Kind     int
		Amount   int
		Wmtranid int
		Purse    string
		Desc     string
	}
}

type IndxOfferMy struct {
	Code  int
	Desc  string
	Value []struct {
		Toolid  int
		Offerid int
		Name    string
		Kind    int
		Price   float32
		Notes   int
		Stamp   int
	}
}

type IndxOfferList struct {
	Code  int
	Desc  string
	Value []struct {
		Offerid int
		Kind    int
		Price   float32
		Notes   int
	}
}

type IndxOfferAdd struct {
	Code  int
	Desc  string
	Value struct {
		Code    int
		OfferID int
	}
}

type IndxOfferDelete struct {
	Code  int
	Desc  string
	Value struct {
		Code    int
		OfferID int
	}
}

type IndxTick struct {
	Code  int
	Desc  string
	Value []struct {
		T     int
		Min   float32
		Max   float32
		Open  float32
		Close float32
		Avg   float32
		Total int
	}
}

// Private methods

// Подсчитывает сигнатуру для запроса
func (i *Indx) signature(str string) (base64String string) {

	// -> sha256
	encoder := sha256.New()
	encoder.Write([]byte(str))
	encodedSum := encoder.Sum(nil)

	// -> base64
	base64String = base64.StdEncoding.EncodeToString(encodedSum)

	return base64String
}

// Выполняет предварительные действия
func (i *Indx) doBeforeActions(fields []string) {
	var sign string = ""
	var size int = len(fields)
	for i := 0; i < size; i++ {
		sign = sign + fields[i]
		if i != (size - 1) {
			sign = sign + ";"
		}
	}
	i.sign = i.signature(sign)
}

// Отправляет запрос на сервер и возвращает ответ сервера
func (i *Indx) sendRequest(request *IndxRequest) (response, sendError string) {
	fmt.Println()
	var url string = URL + request.method
	reader := strings.NewReader(request.data)
	result, err := http.Post(url, "text/json", reader)

	if err != nil { // HTTP response error code
		return "", err.Error()
	}

	body, readAllErr := ioutil.ReadAll(result.Body)

	if readAllErr != nil { // read error code
		return "", readAllErr.Error()
	}
	result.Body.Close()
	return string(body), ""
}

// Public methods

//
//	Balance method. Текущий баланс Трейдера.
//
//	@return	response	IndxBalance	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) Balance() (response IndxBalance, errText string) {

	const method string = "Balance"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)

	if errMarshal != nil {
		return IndxBalance{}, errMarshal.Error()
	}
	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}

	r, errSend := i.sendRequest(&request) // TODO

	if errSend != "" {
		return IndxBalance{},errSend
	}
	
	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxBalance{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	Tools method. Список инструментов биржи.
//
//	@return	response	IndxTools	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) Tools() (response IndxTools, errText string) {

	const method string = "Tools"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxTools{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}

	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxTools{},errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxTools{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	HistoryTrading method. История торгов Трейдера.
//
//	@param	id		string 	Номер инструмента
//	@param	dateStart	string	Начальная дата
//	@param	dateEnd		string	Конечная дата
//
//	@return	response	IndxHistoryTrading	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) HistoryTrading(id, dateStart, dateEnd string) (response IndxHistoryTrading, errText string) {

	const method string = "HistoryTrading"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, id, dateStart, dateEnd})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
		"Trading": map[string]string{
			"ID":        id,
			"DateStart": dateStart,
			"DateEnd":   dateEnd,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxHistoryTrading{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}

	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxHistoryTrading{},errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxHistoryTrading{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	HistoryTransaction method. История трансакций Трейдера.
//
//	@param	id		string	Номер инструмента
//	@param	dateStart	string	Начальная дата
//	@param	dateEnd		string	Конечная дата
//
//	@return	response	IndxHistoryTransaction	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) HistoryTransaction(id, dateStart, dateEnd string) (response IndxHistoryTransaction, errText string) {

	const method string = "HistoryTransaction"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, id, dateStart, dateEnd})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
		"Trading": map[string]string{
			"ID":        id,
			"DateStart": dateStart,
			"DateEnd":   dateEnd,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxHistoryTransaction{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}

	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxHistoryTransaction{}, errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxHistoryTransaction{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	OfferMy method. Список текущих заявок Трейдера на покупку/продажу.
//
//	@return	response	IndxOfferMy	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) OfferMy() (response IndxOfferMy, errText string) {

	const method string = "OfferMy"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxOfferMy{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}
	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxOfferMy{}, errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxOfferMy{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	OfferList method. Список текущих заявок по инструменту на покупку/продажу на бирже.
//
//	@param	id		string	Номер инструмента
//
//	@return	response	IndxOfferList	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) OfferList(id string) (response IndxOfferList, errText string) {

	const method string = "OfferList"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, id})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
		"Trading": map[string]string{
			"ID": id,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxOfferList{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}
	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxOfferList{},errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxOfferList{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	OfferAdd method. Постановка новой заявки Трейдера на покупку/продажу по инструменту на бирже.
//
//	@param	id		string	Номер инструмента
//	@param	count		string	Количество
//	@param	isAnonymous	bool	Статус подачи заявки
//	@param	isBid		bool	Тип подачи заявки
//	@param	price		string	Цена
//
//	@return	response	IndxOfferAdd	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) OfferAdd(id, count string, isAnonymous, isBid bool, price string) (response IndxOfferAdd, errText string) {

	const method string = "OfferAdd"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, id})

	var isAnon string
	if isAnonymous {
		isAnon = "true"
	} else {
		isAnon = "false"
	}
	var isB string
	if isBid {
		isB = "true"
	} else {
		isB = "false"
	}

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
		"Offer": map[string]string{
			"ID":          id,
			"Count":       count,
			"IsAnonymous": isAnon,
			"IsBid":       isB,
			"Price":       price,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxOfferAdd{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}
	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxOfferAdd{}, errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxOfferAdd{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	OfferDelete method. Удаление заявки Трейдера на покупку/продажу инструмента на бирже.
//
//	@param	offerID		string	Номер заявки
//
//	@return	response	IndxOfferDelete	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) OfferDelete(offerID string) (response IndxOfferDelete, errText string) {

	const method string = "OfferDelete"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, offerID})

	type TApiContext struct {
		Login   string
		Wmid    string
		Culture string
		sign    string
	}
	api := TApiContext{
		Login:   i.Login,
		Wmid:    i.Wmid,
		Culture: i.Culture,
		sign:    i.sign,
	}
	type TmpRequest struct {
		ApiContext TApiContext
		OfferID    string
	}
	tRequest := TmpRequest{
		ApiContext: api,
		OfferID:    offerID,
	}

	jsonTmp, errMarshal := json.Marshal(tRequest)
	if errMarshal != nil {
		return IndxOfferDelete{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}
	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxOfferDelete{}, errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxOfferDelete{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

//
//	Tick method. Статистика сделок за период времени.
//
//	@param	tickID		string	Номер инструмента
//	@param	tickKind	string	Период отбора
//
//	@return	response	IndxTick	структура ответа сервера
//	@return errText		string		текст ошибки
//
func (i *Indx) Tick(tickID, tickKind string) (response IndxTick, errText string) {

	const method string = "tick"

	i.doBeforeActions([]string{i.Login, i.Password, i.Culture, i.Wmid, tickID, tickKind})

	var tmp map[string]map[string]string = map[string]map[string]string{
		"ApiContext": map[string]string{
			"Login":     i.Login,
			"Wmid":      i.Wmid,
			"Culture":   i.Culture,
			"Signature": i.sign,
		},
		"Tick": map[string]string{
			"ID":   tickID,
			"Kind": tickKind,
		},
	}

	jsonTmp, errMarshal := json.Marshal(tmp)
	if errMarshal != nil {
		return IndxTick{}, errMarshal.Error()
	}

	var data string = string(jsonTmp)

	var request = IndxRequest{method: method, data: data}
	r, errSend := i.sendRequest(&request)
	if errSend != "" {
		return IndxTick{}, errSend
	}

	errUnmarshal := json.Unmarshal([]byte(r), &response)
	if errUnmarshal != nil {
		return IndxTick{}, errUnmarshal.Error()
	}
	return response, IndxCodes[string(response.Code)]
}

// End of methods
