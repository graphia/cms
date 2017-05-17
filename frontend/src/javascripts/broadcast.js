class CMSMessage {

	constructor(type, alert, content, timeout) {
		this.type = type;		// the alert (Bootstrap) context (success, primary, etc)
		this.alert = alert;		// the bold prefix for the alert
		this.content = content; // the alert's main content
		this.timeout = timeout; // the duration the message will be shown for
		this.active = true; 	// only 'active' alerts are displayed

		// set auto expiry into motion
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

	addMessage(type, alert, content, timeout = 10) {
		console.debug("received broadcast content", type, alert, content)

		let message = new CMSMessage(type, alert, content, timeout);
		this.messages.push(message);

	};

};