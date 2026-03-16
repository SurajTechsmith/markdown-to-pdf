export namespace converter {
	
	export class FolderResult {
	    File: string;
	    Error: string;
	    Duration: string;
	
	    static createFrom(source: any = {}) {
	        return new FolderResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.File = source["File"];
	        this.Error = source["Error"];
	        this.Duration = source["Duration"];
	    }
	}

}

