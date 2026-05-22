export namespace main {
	
	export class LargeFileInfo {
	    path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new LargeFileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.size = source["size"];
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

