export namespace backend {
	
	export class SearchResultItem {
	    id: string;
	    title: string;
	    uploaderName: string;
	    duration: number;
	    durationFormatted: string;
	    thumbnailUrl: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResultItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.uploaderName = source["uploaderName"];
	        this.duration = source["duration"];
	        this.durationFormatted = source["durationFormatted"];
	        this.thumbnailUrl = source["thumbnailUrl"];
	    }
	}
	export class SearchResponse {
	    items: SearchResultItem[];
	    continuationToken?: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], SearchResultItem);
	        this.continuationToken = source["continuationToken"];
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
	
	export class StreamResult {
	    videoId: string;
	    title: string;
	    artist: string;
	    thumbnail: string;
	    audioUrl: string;
	    videoUrl: string;
	    duration: number;
	
	    static createFrom(source: any = {}) {
	        return new StreamResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.videoId = source["videoId"];
	        this.title = source["title"];
	        this.artist = source["artist"];
	        this.thumbnail = source["thumbnail"];
	        this.audioUrl = source["audioUrl"];
	        this.videoUrl = source["videoUrl"];
	        this.duration = source["duration"];
	    }
	}
	export class UserData {
	    favorites: SearchResultItem[];
	    history: SearchResultItem[];
	
	    static createFrom(source: any = {}) {
	        return new UserData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.favorites = this.convertValues(source["favorites"], SearchResultItem);
	        this.history = this.convertValues(source["history"], SearchResultItem);
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

