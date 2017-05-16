class CMSMessage {

	constructor(type, content, timeout) {
		this.type = type;
		this.content = content;
		this.timeout = timeout; // in seconds
		this.active = true;

		this.autoExpire();
	};

	async autoExpire() {
		await this.expireAfter(this.timeout * 1000);
		this.expire();
	}

	// CSS classes to apply to this message
	get classes() {
		return `alert alert-${this.type}`;
	}

	expireAfter(ms) {
		return new Promise(resolve => {
			setTimeout(() => {
				resolve();
			}, ms);
		});
	}

	expire() {
		console.log("expiring...");
		this.active = false;
	}

};

export default class CMSBroadcast {

	constructor() {
		this.messages = [];
	};

	activeMessages() {
		return this.messages.filter(function(message){
			return message.active;
		});
	}

	addMessage(type, content, timeout = 10) {
		console.debug("received broadcast content", type, content)

		let message = new CMSMessage(type, content, timeout);
		this.messages.push(message);

	};

};