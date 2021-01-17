# HTTP API

- [Search](#search)
- [View Page](#view-page)

---

## Search

GET: `/search?q=<query_string>&page=<page_number>`

This API is used for short displaying pages that relevant to query. The result is already sorted from most to least relevant pages.

**Query Params:**

- `q`, String => query for the search.
- `page`, Integer, _OPTIONAL_ => by default set to `1`.

**Example Request:**

```bash
GET /search?q=Hamlet+special
```

**Success Response:**

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
    "ok": true,
    "data": {
        "relevants": [
            {
                "short_html": "<b>Hamlet</b> is my <b>special</b> egg. How are you buddy? I'm not too fond with literature art so I could just...",
                "url": "/pages/123?q=Hamlet+special"
            },
            ...
            {
                "short_html": "This is <b>special</b> drawing created by me. Why don't you just read <b>hamlet</b>?",
                "url": "/pages/743?q=Hamlet+special"
            }
        ],
        "current_page": 1,
        "total_pages": 10
    }
}
```

[Back to Top](#http-api)

---

## View Page

GET: `/pages/<page_number>?q=<query_string>`

This API is used for viewing page in details, basically it's like when you open result from the google book, you will be pointed to the page & your query words will be highlighted.

**Query Params:**

- `q`, String => words that will be highlighted in page

**Example Request:**

```bash
GET /pages/123?q=Hamlet+special
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