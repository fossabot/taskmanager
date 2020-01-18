[![codecov](https://codecov.io/gh/Delgus/taskmanager/branch/master/graph/badge.svg)](https://codecov.io/gh/Delgus/taskmanager)
[![Go Report Card](https://goreportcard.com/badge/github.com/delgus/taskmanager)](https://goreportcard.com/report/github.com/delgus/taskmanager)

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
