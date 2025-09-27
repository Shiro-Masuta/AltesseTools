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

export namespace pkg {
	
	export class FormatStats {
	    count: number;
	    original_size: number;
	    final_size: number;
	
	    static createFrom(source: any = {}) {
	        return new FormatStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.count = source["count"];
	        this.original_size = source["original_size"];
	        this.final_size = source["final_size"];
	    }
	}

}

export namespace services {
	
	export class WidgetStats {
	    total_converted: number;
	    formats: Record<string, pkg.FormatStats>;
	    total_saved_cd: number;
	    total_saved_floppy: number;
	    total_original_size: number;
	    total_compressed_size: number;
	
	    static createFrom(source: any = {}) {
	        return new WidgetStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_converted = source["total_converted"];
	        this.formats = this.convertValues(source["formats"], pkg.FormatStats, true);
	        this.total_saved_cd = source["total_saved_cd"];
	        this.total_saved_floppy = source["total_saved_floppy"];
	        this.total_original_size = source["total_original_size"];
	        this.total_compressed_size = source["total_compressed_size"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace system {
	
	export class CPUMinimalInfo {
	    name: string;
	    cores: number;
	    threads: number;
	    mhz: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUMinimalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.cores = source["cores"];
	        this.threads = source["threads"];
	        this.mhz = source["mhz"];
	        this.percent = source["percent"];
	    }
	}
	export class DiskMinimalInfo {
	    mountpoint: string;
	    total: number;
	    used: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskMinimalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mountpoint = source["mountpoint"];
	        this.total = source["total"];
	        this.used = source["used"];
	        this.percent = source["percent"];
	    }
	}
	export class MemoryMinimalInfo {
	    total: number;
	    used: number;
	    free: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new MemoryMinimalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.used = source["used"];
	        this.free = source["free"];
	        this.percent = source["percent"];
	    }
	}
	export class SystemMinimalInfo {
	    cpu?: CPUMinimalInfo;
	    memory?: MemoryMinimalInfo;
	    disk: DiskMinimalInfo[];
	
	    static createFrom(source: any = {}) {
	        return new SystemMinimalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpu = this.convertValues(source["cpu"], CPUMinimalInfo);
	        this.memory = this.convertValues(source["memory"], MemoryMinimalInfo);
	        this.disk = this.convertValues(source["disk"], DiskMinimalInfo);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

