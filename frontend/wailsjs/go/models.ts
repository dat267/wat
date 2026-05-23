export namespace main {
	
	export class AppConfig {
	    apiKeys: Record<string, string>;
	    settings: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new AppConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.apiKeys = source["apiKeys"];
	        this.settings = source["settings"];
	    }
	}
	export class DiskEntry {
	    name: string;
	    path: string;
	    size: number;
	    isDir: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DiskEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	        this.isDir = source["isDir"];
	    }
	}
	export class SystemStats {
	    cpuPercent: number;
	    memoryPercent: number;
	    memoryUsed: number;
	    memoryTotal: number;
	    diskPercent: number;
	    diskUsed: number;
	    diskTotal: number;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cpuPercent = source["cpuPercent"];
	        this.memoryPercent = source["memoryPercent"];
	        this.memoryUsed = source["memoryUsed"];
	        this.memoryTotal = source["memoryTotal"];
	        this.diskPercent = source["diskPercent"];
	        this.diskUsed = source["diskUsed"];
	        this.diskTotal = source["diskTotal"];
	        this.uptime = source["uptime"];
	    }
	}

}

