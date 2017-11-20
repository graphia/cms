<template>
	<div>
		<h1>User Settings</h1>


		<form @submit="create">
			<label for="ssh-key"></label>

			<textarea
				name="ssh-key"
				v-model="key.raw"
			/>

			<input type="submit">
		</form>
	</div>
</template>

<script>
	import checkResponse from "../../javascripts/response.js";
	import {CMSNewPublicKey,CMSPublicKey} from "../../javascripts/models/public_key.js";

	export default {
		name: "UserSettings",
		data() {
			return {
				key: new CMSPublicKey
			};
		},
		methods: {
			async create(event) {

				console.debug("clicked");
				event.preventDefault();

				let response = await this.key.create();

				if (!checkResponse(response.status)) {
					throw "Could not create key"
				}

				console.log("Key created")
			}
		}
	};
</script>
