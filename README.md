# VariationSeriesVisual_V2
Version_2. The input is a variational series, a set of numbers that can be repeated. The program processes it, produces a histogram, graph and table in Excel format.

В версии 2 проект аккуратно разделен на файлы по архитектуре кода, сделан более удобочитаемым и менее громоздким, добавлены горутины для быстродействия. Добавлены табличные тесты. 

На вход дается вариационный ряд, набор чисел, которые могут повторяться. Программа обрабатывает их и преобразует в таблицу Excel.
Столбцы: вариант(число), количество ее повторения (частота) и относительная частота варианта. 
Затем программа проверяет максимальные и минимальные значения на то, нужно ли их исключить из ряда, рассчитывая критерий и сравнивая его с табличным значением. 
Наконец программа создает гистограмму и полигон распределения в html формате, которые вы можете сохранить.

In version 2, the project is neatly divided into files by code architecture, made more readable and less cumbersome, goroutines were added for speed. Table tests were added too.

The input is a variational series, a set of numbers that can be repeated. The program processes and converts them into an Excel spreadsheet. 
Columns: variant (number), number of repetitions (frequency) and relative frequency of the variant. 
The program then checks the maximum and minimum values to see if they should be excluded from the series by calculating the criterion
and comparing it with the table value. Finally, the program creates a histogram and distribution polygon in html format, which you can save.
