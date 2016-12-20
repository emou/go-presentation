# Уъркшоп към семинара "Програмиране с Go"

Този документ съдържа инструкциите за уъркшопа и вашето участие в него.

# Камък-ножица-хартия

Проектът, с който ще си играем в уъркшопа е имплементация на добре познатата на всички игра
"Камък-ножица-хартия" като мрежова игра.

Играта има 2 основни компонента:

- Сървър, който управлява играта и ние сме написали.
- Клиент, който ще трябва да напишете вие и да играете помежду си, като се свържете към сървъра!

# Инструкции

Тук ще опишем всички необходими стъпки, за да напишете първия си Go код. Не се колебайте да
задавате въпроси или да посочвате неточности.

## Инсталиране на Go

За първоначланата Go следвайте инструкциите на официалния сайт на Go.

- За Windows Пробвайте MSI инсталатора: https://golang.org/doc/install#windows
- За UNIX-базирана ОС, използвайте пакетния си мениджър или метода описан тук: https://golang.org/doc/install#tarball

За да проверите инсталацията си, следвайте инструкциите тук:
https://golang.org/doc/install#testing

## Задача: клиент за играта

Вашата задача ще е да напишете клиент за сървъра, написан на Golang. За да направите това, ще вие е
нужно да знаете как да си говорите със сървъра.

### Инсталирайте кода на Камък-Ножица-Хартия сървъра

```
$ go get github.com/emou/go-presentation/gorps
```

### main.go

Нека създадем `main.go`, където ще напишем кода на клиента със съдържание празна функция:

```
package main

func main() {
}
```

Да проверим, че се билдва и работи:

```
go run main.go
```

### Свързване към сървъра

Ще използваме пакета от стандартната библиотека [net](https://golang.org/pkg/net/), за да
се свържем към сървъра.

```
func main() {
  conn, err := net.Dial("tcp", "golang.org:80")
  if err != nil {
    // handle error
  }
  fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
}
```

### Протокол

Сървърът на играта използва TCP мрежовия протокол и дефинира няколко прости съощения, определящи
"приложния протокол".

Съобщенията, които пращате на съвръва са разделени с нов ред ('\n').

## Бот

Следващата ви задача ще е да напишете бот, който играе камък-ножица-хартия вместо вас. Нямате
ограничение в броя на връзки, които да използвате с едно и също име.