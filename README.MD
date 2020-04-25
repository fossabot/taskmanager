[![Sourcegraph](https://sourcegraph.com/github.com/Delgus/taskmanager/-/badge.svg)](https://sourcegraph.com/github.com/delgus/taskmanager?badge)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/Delgus/taskmanager)
[![codecov](https://codecov.io/gh/Delgus/taskmanager/branch/master/graph/badge.svg)](https://codecov.io/gh/Delgus/taskmanager)
[![Go Report Card](https://goreportcard.com/badge/github.com/delgus/taskmanager)](https://goreportcard.com/report/github.com/delgus/taskmanager)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/delgus/taskmanager/master/LICENSE)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FDelgus%2Ftaskmanager.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2FDelgus%2Ftaskmanager?ref=badge_shield)


### TASK MANAGER (Решение тестовой задачи)

#### Задача:  

Реализовать систему очередей задач (библиотекой).  
Разруливать задачи должен TaskManager, задачи могут быть различных типов.  
У задач могут быть различные уровни приоритетов:

- 1 - очень низкий
- 2 - низкий
- 3 - средний
- 4 - высокий
- 5 - очень высокий

У задач есть события:

- Создание
- Выполнение
- Завершение
- Ошибка выполнения

На любое событие можно вешать от одного и более обработчиков, они должны выполняться по очередности FIFO.

Наличие тестов обязательное, и документирование в стиле godoc плюсом. Использовать готовые решения и библиотеки 
запрещено.


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2FDelgus%2Ftaskmanager.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2FDelgus%2Ftaskmanager?ref=badge_large)