# Hraaska
## Url shortner service in golang using mongodb as a database
<p>&nbsp;</p>

# APIs

## 1. Create Short URL 
</br>

### REQUEST  
</br>

* **URL** - `/api/url/shorten`  
* **Method** - `POST`  
* **Body** -  
```
    {
        "url":"https://google.com"
    }
```

### RESPONSE  
</br>

* **Success Response**  
    **Code** - `200`  
    **Response Body** -  
```
    {
        "status":200,
        "success":true,
        "data":{
            "id":"id_of_inserted_record",
            "url":"https://google.com",
            "hash":"acacwwf"
        }
    }
```
* **Error Response**  
    **Code** - ```4xx|5xx ```  
    **Response body** -  
```
    {
        "status":"4xx|5xx"
        "success":false,
        "error":{
            "message":"error_message"
        }
    }
```

---

## 2. Redirect to destination  
</br>

## REQUEST  
</br>  

* **URL** - `/api/url/redirect`  
* **Method** - `POST`  
* **Body** -  
```
    {
        hash: "acacwwf"
    }
```  
## RESPONSE  
</br>

* **Success Response**  
    **Code** - `301`

* **Error Response**  
    **Code** - `4xx|5xx`  
    **Response Body** - 
```
    {
        "status":"4xx|5xx"
        "success":false,
        "error":{
            "message":"error_message"
        }
    }
```  
