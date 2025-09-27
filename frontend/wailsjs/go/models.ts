export namespace images {
	
	export class Options {
	    Format: string;
	    Quality: number;
	    Lossless: boolean;
	    PNGLevel: number;
	    TIFFCompress: number;
	
	    static createFrom(source: any = {}) {
	        return new Options(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Format = source["Format"];
	        this.Quality = source["Quality"];
	        this.Lossless = source["Lossless"];
	        this.PNGLevel = source["PNGLevel"];
	        this.TIFFCompress = source["TIFFCompress"];
	    }
	}

}

export namespace main {
	
	export class FileData {
	    name: string;
	    content: number[];
	
	    static createFrom(source: any = {}) {
	        return new FileData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.content = source["content"];
	    }
	}

}

