# textextract

textextract is a tiny library (87 lines of Go) that identifies where the article content is in a HTML page (as opposed to navigation, headers, footers, ads, etc), extracts it and returns it as a string.

It's a tree search and score algorithm, it uses a very simple scoring rule, it is surprisingly effective.

## What it's for

If you're doing semantic analysis on crawled information and you need article content to feed into some other process like a semantic extractor, classifier, etc. It preserves the rendering order of text but it doesn't preserve white space.

## How it works

1. It parses the HTML into a node tree using the standard Go html package.

2. It then walks the tree depth first and scores each node on route. The score from the parent node is pushed down as the basis for the child node. The scoring formula is: WORDCOUNT - WORDCOUNTINANCHOR^2, where WORDCOUNT is the number of words in the node that are not hyperlinked and WORDCOUNTINANCHOR is the number of words in the node that are hyperlinked. The WORDCOUNTINANCHOR for each node is actually calculated as 1 + WORDCOUNTINACHOR^2 just because there is often anchors on things other than words so WORDCOUNTINACHOR^2 is often zero.

3. As it goes it'll add nodes that are below the minimum score to a toDelete slice. When the recursion is finished, it'll delete all nodes in the toDelete slice.

4. Finally, the filtered tree is parsed again, depth first, and all text nodes are printed to a string.

## How to install it

    go get https://github.com/emiruz/textextract

## How to use it

    import "github.com/emiruz/textextract"

    main func() {
    	textextract.MinScore = 5 // the default is 5.
        extractedText, err := textextract.ExtractFromHtml(yourUTF8HTMLString)
    }

## License

MIT Licensed, do as you will with it.

## Bugs

Please submit them as issues on the repository.

## TODO

1. Add tests

2. Add comments.
