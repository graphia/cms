<template>
	<div class="col-md-12">
		<h1>User Settings</h1>





		<div id="existing-keys">
			<h4>Existing keys</h4>
			<ul class="list-group">

				<li :id="`ssh-public-key-${key.id}`" class="list-group-item d-flex justify-content-between align-items-center" v-for="(key, i) in keys" :key="i">

					<h3>{{ key.name }}</h3>

					<code class="key-data">
						{{ key.fingerprint }}
					</code>

					<button class="btn btn-danger delete-pk-button" :data-key-id="key.id" @click="deleteKey">
						Delete
					</button>
				</li>
			</ul>
		</div>


		<div id="new-ssh-key" class="mt-4">
			<h4>Upload a new SSH key</h4>

			<form @submit="create">

				<div class="form-group">
					<label for="name">Name</label>
					<input
						type="text"
						v-model="newKey.name"
						class="form-control"
						placeholder="laptop"
						required
					/>
				</div>

				<div class="form-group">
					<label for="ssh-key">SSH Key</label>

					<textarea
						name="ssh-key"
						v-model="newKey.raw"
						class="form-control"
						placeholder="ssh-rsa ABC123…"
					/>

				</div>

				<div class="form-group">
					<input class="form-control btn btn-success" type="submit" value="Create SSH Key">
				</div>
			</form>
		</div>
	</div>
</template>

<script>
	import checkResponse from "../../javascripts/response.js";
	import CMSPublicKey from "../../javascripts/models/public_key.js";

	export default {
		name: "UserSettings",
		data() {
			return {
				newKey: new CMSPublicKey,
				keys: []
			};
		},
		created() {
			this.refresh();
		},
		methods: {
			async create(event) {

				console.debug("clicked");
				event.preventDefault();

				let response = await this.newKey.create();

				if (!checkResponse(response.status)) {
					throw "Could not create key";
				};

				this.refresh();
				this.reset();

				this.$store.state.broadcast.addMessage(
					"success",
					"SSH Key Created",
					"You can now clone the content to your computer and work offline",
					3
				);

			},
			async refresh() {

				try {
					let response = await CMSPublicKey.all();

					if (!checkResponse(response.status)) {
						throw "Could not fetch public keys";
					};

					let json = await response.json();

					let keys = json.map((key) => {
						return new CMSPublicKey(key);
					});

					this.keys = keys;

				}
				catch(e) {
					console.error(e);
				}

			},
			async deleteKey(event) {
				event.preventDefault();

				try {
					let id = event.currentTarget.getAttribute('data-key-id')

					let marked = this.keys.find((key) => { return key.id == id });

					if (!marked) {
						console.error(`key ${id} not found in`, keys);
							throw {message: "no matching key present", response: response}
					};

					let response = await marked.delete();

					if (!checkResponse(response.status)) {
						throw {message: "deletion failed", response: response}
					};

					// key deleted
					this.$store.state.broadcast.addMessage(
						"warning",
						"SSH Key Deleted",
						`You can no longer use key '${marked.name}' for SSH`,
						3
					);

					this.refresh();

				}
				catch(e) {
					console.error("Could not delete key", e)
				}

			},
			reset() {
				this.newKey = new CMSPublicKey;
			}
		}
	};
</script>

<style lang="scss">
	#existing-keys {

		li {

			pre {
				white-space: pre;
			};

			.delete-pk-button {
				text-align: right;
			};

		};

	};
</style>
