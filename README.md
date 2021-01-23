# ShakeSearch

I'm not a fan of literary arts even in my own native language. So I'm not so sure how to approach this problem correctly.

But when I was writing my final thesis in college, I need to read many books for my research. So I think I could somehow use this experience to approach this issue.

## Basic Idea

When I want to search about something, essentially I want to know more about the context surrounding that something. For example when I want to search about `canned answer`, what I was expecting to get from my search is anything that related to `canned answer`. For example how to implement it in chatbot, does it performs well in real world, who is people who have been trying to apply this method in real life, etc...

So I was thinking maybe all of these stuff is also the same for people who are seeking information about literary arts. I mean when people search something about literary arts, maybe they want to know more anything that related to that art. For example, when people search about `Hamlet`, maybe they want to know more about:

- what is hamlet
- why the name is hamlet
- what the story that revolves around hamlet
- etc...

So I was thinking, why don't we just try to build index which relate every keywords with documents that related to that keywords? Just like Lucene Index but in our very own way?

I mean the source file contains complete works of Shakespeare, right? That means naturally there are multiples works inside it. If we can break these works somehow into document like page, maybe we could build an index where we could map every keywords with specific documents that related to them, right?

By having this index, I think we could help people to find contexts related to the keywords they typed in. If they want to find out more about the context, they could just open the full page and read more about it.

## How It's Work

So basically my approach for solving this problem is like this:

1. Break the source file into pages
2. Associate these pages with every keywords inside them in an index
3. When user type in the keywords, we will give them associated documents / pages, sorted by relevance
4. When user want to explore more the contexts, they just need to open the full document and move to the next or previous documents as necessary

## Link to Demo

https://stormy-temple-09337.herokuapp.com/

## Future Plan Ideas

1. Improve the formula for scoring relevant documents, current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/index/index.go#L112-L127). The issuw with current formula is currently it is not yet involving words ordering. I believe if we are able to do this, we could improve the search result quality greatly. Because usually when people put more words in search query, they are trying to search for something more specific.

2. Improve document preview on search result. The current one is ok (?), but I think it should be much better. The current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/doc/doc.go#L63-L78).

3. Improve the algorithm to extract the keywords. Current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/index/index.go#L33-L72).

4. Ask users feedback by letting them use the product, usually we will find a lot of insight from that. 😁

5. Handling misspelling?