var app = new Vue({
    el: '#app',
    data: {
        search_query: "",
        relevants: [],
        current_page: 1,
        prev_page: null,
        next_page: null,
        current_doc_page: null,
        prev_doc_page: null,
        next_doc_page: null,
        doc_found_words: [],
        doc_data: ""
    },
    methods: {
        async search() {
            try {
                // fetch search result
                const response = await fetch(`/search?q=${this.search_query}&page=${this.current_page}`);
                const result = await response.json();
                // set state
                this.relevants = result.data.relevants;
                this.prev_page = result.data.prev_page;
                this.next_page = result.data.next_page;
                // alert if no result
                if (!this.relevants.length || this.relevants.length == 0) {
                    alert("No Result");
                    return
                }
            } catch($e) {
                console.log($e);
            }
        },
        async initSearch() {
            this.current_page = 1;
            await this.search();
        },
        async nextSearch() {
            this.current_page = this.next_page;
            await this.search();
        },
        async prevSearch() {
            this.current_page = this.prev_page;
            await this.search();
        },
        async openDocPage(id, foundWords) {
            try {
                // fetch page data
                const response = await fetch(`/pages/${id}?q=${foundWords.join(",")}`);
                const result = await response.json();
                // set doc state
                this.current_doc_page = result.data.current_page;
                this.next_doc_page = result.data.next_page;
                this.prev_doc_page = result.data.prev_page;
                this.doc_found_words = foundWords;
                this.doc_data = result.data.body_html;
                // open modal
                $('#modal').addClass('is-active');
            } catch($e) {
                console.log($e);
            }
        },
    }
});