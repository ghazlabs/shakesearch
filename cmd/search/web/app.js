var app = new Vue({
    el: '#app',
    data: {
        search_query: "",
        relevants: [],
        current_page: 1,
        prev_page: null,
        next_page: null
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
        }
    }
});