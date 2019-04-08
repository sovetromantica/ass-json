# ASS->JSON->ASS
This program is a converter of ASS/SSA to JSON and vice versa
## Syntax
At the moment it has 3 arguments:
```
-h --help
-a --ass file.ass -- Convert .ass in to JSON-array (only stdout yet)
-j --json file.json -- Convert .json in to ASS (same stdout only)
```
 
Known problems:
* Handling of fonts doesn't work yet

# ASS->JSON->ASS
Данная программа представляет собой парсер формата ASS/SSA в JSON и обратно.

В текущий момент имеет всего 3 аргумента:
```
-h --help -- выводит информацию о командах
-a --ass file.ass - обработает ASS файл и выведет в stdout JSON массив
-j --json file.json - обработает JSON файл и выведет ASS в stdout
```

Известные проблемы:
* Не умеет работать со шрифтами.
