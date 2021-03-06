# HTTP API

This document contains specifications for the APIs that will be used in this project.

Table of contents:

- [Search](#search)
- [View Page](#view-page)

---

## Search

GET: `/search?q=<query_string>&page=<page_number>`

This API is used for short displaying pages that relevant to query. The result is already sorted from most to least relevant pages.

Page number started from `1`.

**Query Params:**

- `q`, String => query for the search, url encoded.
- `page`, Integer, _OPTIONAL_ => by default set to `1`.

**Example Request:**

```bash
GET /search?q=Hamlet+special
```

**Success Response:**

- Normal Response:

    ```json
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "ok": true,
        "data": {
            "relevants": [
                {
                    "id": 123,
                    "title": "Page 123",
                    "short_html": "<b>Hamlet</b> is my <b>special</b> egg. How are you buddy? I'm not too fond with literature art so I could just...",
                    "found_words": ["hamlet", "special"],
                    "score": 89.0,
                    "url": "/pages/123?q=hamlet,special"
                },
                {
                    "id": 743,
                    "title": "Page 743",
                    "short_html": "This is <b>special</b> drawing created by me. Why don't you just read <b>hamlet</b>?",
                    "found_words": ["hamlet", "special"],
                    "score": 70.0,
                    "url": "/pages/743?q=hamlet,special"
                }
            ],
            "current_page": 1,
            "prev_page": null,
            "next_page": 1,
            "total_pages": 10
        }
    }
    ```

- No Result Response:

    ```json
    HTTP/1.1 200 OK
    Content-Type: application/json

    {
        "ok": true,
        "data": {
            "relevants": [],
            "current_page": 1,
            "total_pages": 1
        }
    }
    ```


**Error Responses:**

- Empty Query

    ```json
    HTTP/1.1 400 Bad Request
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_EMPTY_QUERY"
    }
    ```

    Client will receive this error when it is sending empty query to server.

- Page Not Found

    ```json
    HTTP/1.1 404 Not Found
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_PAGE_NOT_FOUND"
    }
    ```

    Client will receive this error when it is trying to access non existent page. For example the total pages is `10`, but the client tried to access the page `11`.

[Back to Top](#http-api)

---

## View Page

GET: `/pages/<page_number>?q=<query_string>`

This API is used for viewing page in details, basically it's like when you open result from the google book, you will be pointed to the page & your query words will be highlighted.

Page number started from `1`.

**Query Params:**

- `q`, String, _OPTIONAL_ => words that will be highlighted in the page separated by comma

**Example Request:**

```bash
GET /pages/123?q=Hamlet,special
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "body_html": "<span style=\"highlight\">Hamlet</span> is my <span style=\"highlight\">special</span> egg. How are you buddy? I'm not too fond with literature art so I could just try to work on it.",
        "current_page": 123,
        "prev_page": 122,
        "next_page": 124,
        "total_pages": 10000
    }
}
```

**Error Responses:**

- Page Not Found

    ```json
    HTTP/1.1 404 Not Found
    Content-Type: application/json

    {
        "ok": false,
        "err": "ERR_PAGE_NOT_FOUND"
    }
    ```

    Client will receive this error when it is trying to access non existent page. For example the total pages is `10`, but the client tried to access the page `11`.

[Back to Top](#http-api)

---