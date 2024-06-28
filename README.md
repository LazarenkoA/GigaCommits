# GigaCommits
### Описание
GigaCommits - это CLI инструмент для генерации сообщений для git commit на основе ИИ, использующий GigaChat API от Сбера. Этот инструмент помогает создавать информативные и осмысленные сообщения для коммитов, облегчая процесс ведения истории изменений в вашем проекте.



**Предварительные требования**<br>
Перед началом работы убедитесь, что у вас установлены следующие инструменты:

* Git
* Go (версия 1.21 и выше)
* Доступ к GigaChat API (необходим client ID и client Secret)

### Установка
```shell
go install https://github.com/LazarenkoA/GigaCommits
```

### Использование
```shell
gigacommit
```


<sub><sup>Проект был вдохновлен аналогичным инструментом [aicommits](https://github.com/Nutlope/aicommits)<br>
Проект работает на базе [gigachat api](https://developers.sber.ru/docs/ru/gigachat/api/reference/rest/post-chat) с использованием клиента [Go GigaChat](https://github.com/paulrzcz/go-gigachat)</sub></sup>