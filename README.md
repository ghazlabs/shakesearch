# ShakeSearch

I'm not a fan of literary arts even in my own native language. So I'm not so sure how to approach this problem correctly.

But when I was writing my final thesis in college, I need to read many books for my research (fyi it was basically research about chatbot but for voice). So I think I could somehow use this experience to approach this issue.

When I want to search about something, essentially I want to know more about the context surrounding that something. For example when I want to search about `canned answer`, what I was expecting to get from my search is anything that related to `canned answer`. For example how to implement it, how its performance in real world, people that have been using this method to solve their challenge, etc...

## Basic Idea

So I was thinking maybe all of these also the same in literary arts. I mean when people search something about literary arts, maybe they want to know more anything that related to that. For example, when people search about `Hamlet`, maybe they want to know more about:

- what is hamlet
- why the name is hamlet
- what the story that revolves around hamlet
- etc...

So I was thinking, why don't we just try to build index which relate every keywords with documents that related to that keywords? Just like Lucene Index but in our very own way?

I mean the source file contains complete works of Shakespeare, right? That means naturally there are multiples works inside it. If we can break these works somehow into document like page, maybe we could build an index where we could map every keywords with specific documents that related to them, right?

By having this index, I think we could help people to find contexts related to the keywords they typed in. If they want to find out more about the context, they could just open the full page and read more about it.

## Link to Demo

https://stormy-temple-09337.herokuapp.com/

## How It's Work

So basically my approach for solving this problem is like this:

1. Break the source file into pages
2. Associate these pages with every keywords inside them in an index
3. When user type in the keywords, we will give them associated documents / pages, sorted by relevance
4. When user want to explore more the contexts, they just need to open the full document and move to the next or previous documents as necessary

## Future Plan Ideas

1. Improve the formula for scoring relevant documents, current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/index/index.go#L112-L127). Notice that current formula is not yet counting the word ordering which I believe will improve the search result quality. Because usually when people put more words in search query, they are trying to search for something more specific.

2. Improve document preview on search result. The current one is ok (?), but I think it should be much better. The current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/doc/doc.go#L63-L78).

3. Improve the algorithm to extract the keywords. Current implementation could be seen [here](https://github.com/ghazlabs/shakesearch/blob/master/internal/index/index.go#L33-L72).

4. Ask users feedback by letting them use the product, usually we will find a lot of insight from that. 😁

5. Handling misspelling?