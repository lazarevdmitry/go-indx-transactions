# Go Indx Transactions

Библиотека, реализующая взаимодействие с API Интернет-биржи INDX.

При помощи данной библиотеки вы сможете создавать собственные программы, работающие на Интернет-бирже INDX.

## Как использовать эту библиотеку?


### Скачиваем

1. Создайте пустую папку.
2. Войдите в эту папку через консоль/терминал и выполните команду (подразумевается, что у вас уже установлен git):
```bash
   git clone https://github.com/lazarevdmitry/go-indx-transactions.git
```
   Или скачайте файл go-indx-transactions-master.zip и распакуйте его на своем жестком диске.

### Компилируем

1. В папке go-indx-transactions выполните команды (у вас должна быть предварительно установлена переменная среды GOPATH):
```bash
   go build
```
   Затем
```bash
   go install
```

### Применяем

Пример:

```go
// файл main.go
package main

import indx

func main(){
     // ...какой-то код (например, объявление и инициализации переменных для запроса)

     // Первый вариант инициализации структуры Indx
     exchange := indx.Indx{} // Создает пустую структуру типа Indx
     exchange.Login = "john"
     exchange.Password = "qwerty"
     exchange.Wmid = "1234567890"
     exchange.Culture = "ru-RU"
     
     // Второй вариант инициализации структуры Indx
     exchange = Indx{
     	Login: "john",
        Password: "qwerty",
        Wmid: "1234567890",
        Culture: "ru-RU",
     }

     // Отправка запросов (каждый из методов возвращает свою структуру)
     balance,_ := exchange.Balance()							// возвращает IndxBalance{}
     tools,_ := exchange.Tools()	     							// возвращает IndxTools{}
     historyTrading,_ := exchange.HistoryTrading("100500","20190101","20191231")		// возвращает HistoryTrading{}
     historyTransaction,_ := exchange.HistoryTransaction("100500","20190101","20191231")	// возвращает HistoryTransaction{}
     offerMy,_ := exchange.OfferMy()							// возвращает OfferMy{}
     offerList,_ := exchange.OfferList("1")						// возвращает OfferList{}
     offerAdd,_ := exchange.OfferAdd("1","120",true,true,"1000.00")			// возвращает OfferAdd{}
     offerDel,_ := exchange.OfferDelete("100500")						// возвращает OfferDelete{}
     statistics,_ := exchange.Tick("1","4")						// возвращает Tick{}
     // ... далее какие-то действия ...
}
```

## Предупреждение

Внимание! Методы HistoryTrading, HistoryTransaction,OfferMy,OfferList,OfferAdd,OfferDelete,Tick пока что не тестировались.

