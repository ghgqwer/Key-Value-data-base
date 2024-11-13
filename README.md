# Bolshoi Golang Project

## 📖 Описание

В рамках первого проекта вам требуется сделать **in-memory key-value базу данных** с **HTTP интерфейсом**, возможностью **персистентного хранения** и **автоматического удаления устаревших записей**.

## ✨ Типы хранимых данных

Все ключи базы данных - **строки**. По одному ключу может быть только одно значение одного из типов:

- **Скаляр**
- **Словарь**
- **Массивы**

### 🏷️ Скаляр

Скаляр - единичное значение типа строка либо целое число. Способ внутреннего представления скаляров - на ваш выбор (можно хранить все как строки, как интерфейсы и т.д.). При чтении данных у пользователя должна быть возможность явно узнать тип скаляра (строка/число).

#### ⚙️ Операции по работе со скалярами

##### GET key

Возвращает значение по ключу key. Если по ключу нет значения, возвращается nil. Если значение по ключу не типа скаляр, возвращается ошибка. Возвращаемые строки должны быть заключены в кавычки, возвращаемые числа должны быть без кавычек.

**Пример:**
```plaintext
SET name "Anton"
OK
SET number 42
OK
GET name
"Anton"
GET PUE
(nil)
GET number
42
```
##### GET key
#####SET key value [EX seconds]

Устанавливает значение по ключу key равным value. Если value указано в кавычках - это строка, иначе это число.

**Пример:**
```plaintext

SET key1 "value1" EX 20
OK
SET key2 value1 EX 20
Requested value is not a number
```
