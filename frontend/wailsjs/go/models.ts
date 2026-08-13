export namespace registry {
	
	export class Tool {
	    id: string;
	    name: string;
	    description: string;
	    category: string;
	    keywords: string[];
	
	    static createFrom(source: any = {}) {
	        return new Tool(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.keywords = source["keywords"];
	    }
	}

}

