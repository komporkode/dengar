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
	export class Playlist {
	    id: string;
	    name: string;
	    description: string;
	    tracks: SearchResultItem[];
	    createdAt: number;
	    updatedAt: number;
	
	    static createFrom(source: any = {}) {
	        return new Playlist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.tracks = this.convertValues(source["tracks"], SearchResultItem);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
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
	    bitrate?: number;
	    format?: string;
	    quality?: string;
	
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
	        this.bitrate = source["bitrate"];
	        this.format = source["format"];
	        this.quality = source["quality"];
	    }
	}
	export class UserData {
	    favorites: SearchResultItem[];
	    history: SearchResultItem[];
	    playlists: Playlist[];
	    playCounts: {[key: string]: number};
	    trackMeta?: {[key: string]: SearchResultItem};
	
	    static createFrom(source: any = {}) {
	        return new UserData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.favorites = this.convertValues(source["favorites"], SearchResultItem);
	        this.history = this.convertValues(source["history"], SearchResultItem);
	        this.playlists = this.convertValues(source["playlists"], Playlist);
	        this.playCounts = source["playCounts"];
	        this.trackMeta = this.convertValues(source["trackMeta"], SearchResultItem, true);
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

export namespace main {
	
	export class HomeData {
	    topPlayed: backend.SearchResultItem[];
	    recommended: backend.SearchResultItem[];
	
	    static createFrom(source: any = {}) {
	        return new HomeData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.topPlayed = this.convertValues(source["topPlayed"], backend.SearchResultItem);
	        this.recommended = this.convertValues(source["recommended"], backend.SearchResultItem);
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
	export class UpdateAsset {
	    name: string;
	    downloadUrl: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateAsset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.downloadUrl = source["downloadUrl"];
	        this.size = source["size"];
	    }
	}
	export class UpdateInfo {
	    currentVersion: string;
	    latestVersion: string;
	    hasUpdate: boolean;
	    releaseUrl: string;
	    releaseNotes: string;
	    publishedAt: string;
	    downloadUrl: string;
	    assetFileName: string;
	    assetSize: number;
	    assets: UpdateAsset[];
	
	    static createFrom(source: any = {}) {
	        return new UpdateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentVersion = source["currentVersion"];
	        this.latestVersion = source["latestVersion"];
	        this.hasUpdate = source["hasUpdate"];
	        this.releaseUrl = source["releaseUrl"];
	        this.releaseNotes = source["releaseNotes"];
	        this.publishedAt = source["publishedAt"];
	        this.downloadUrl = source["downloadUrl"];
	        this.assetFileName = source["assetFileName"];
	        this.assetSize = source["assetSize"];
	        this.assets = this.convertValues(source["assets"], UpdateAsset);
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

