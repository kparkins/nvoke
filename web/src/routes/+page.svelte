<script>
	import { onMount } from "svelte";

	let query = "";
	let result = "";  
	onMount(() => {
	});

	async function handleSearch() {
	  if (query.trim()) {
		// sound.currentTime = 0;
		// sound.play();
		// Here, we stub out an API call
		const response = await fetch('http://api.local/v1/completion', {
		  method: 'POST',
		  headers: { 'Content-Type': 'application/json' },
		  body: JSON.stringify({ 'query': query })
		});
  
		const data = await response.json();
		console.log(data.answer); // Log the response for now
		result = data.answer.replace(/\n/g, '<br>');; 
	  }
	}

	  // background-image: url('heavenly-background.jpg');
  </script>
  
  <style>
	* {
		margin: 0px;
		padding: 0px;
	}
/* 
	body {
	  font-family: 'Times New Roman', serif;
	  background-size: cover;
	  background-position: center;
	  color: #fff;
	  margin: 0;
	  height: 100vh;
	  display: flex;
	  align-items: center;
	  justify-content: center;
	} */
  
	.search-container {
	  display: flex;
	  flex-direction: row;
	  align-items: center;
	  justify-content: center;
	  width: auto;
	  height: auto;
	  max-width: 800px;
	  padding: 10px 20px;
	  background-color: rgba(255, 255, 255, 0.2);
	  border-radius: 10px;
	  box-shadow: 1px 4px 10px rgba(0, 0, 0, 0.5);
	  margin: 10% 20%;
	  margin-top: 30%;
	}
  
	input[type="text"] {
	  padding: 10px;
	  font-size: 16px;
	  border: 2px solid white;
	  border-radius: 5px;
	  outline: none;
	  color: #000;
	  background-color: #fff;
	  width: 100%;
	  height: 100%;
	}
  
	button {
	  padding: 10px 20px;
	  font-size: 16px;
	  background-color: hwb(223 31% 0%);
	  color: #ffffff;
	  border: 2px solid white;
	  border-radius: 5px;
	  cursor: pointer;
	  transition: background-color 0.3s, color 0.3s;
	  width: 20%;
	  max-width: 150px;
	}
  
	button:hover {
	  background-color: hwb(223 52% 2%);
	  color: hwb(0 100% 0%);
	}
  
	@media (min-width: 600px) {
	  .search-container {
		flex-direction: row;
	  }
  
	  input[type="text"] {
		margin-bottom: 0;
		margin-right: 10px;
	  }
	}

		/* New styles for the result display */
	.result-container {
		display: flex;
		padding: 10px 20px;
		margin: 10% 20%;
		align-items: center;
		justify-content: center;
		width: 80%;
		max-width: 800px;
		text-wrap: wrap;
		background-color: rgba(255, 255, 255, 0.2);
		border-radius: 10px;
		box-shadow: 1px 4px 10px rgba(0, 0, 0, 0.5);
		width: auto;
		height: auto;
		min-height: 40%;
	}
	.empty-result-container {
		display: flex;
		padding: 25% 20px;
		margin: 10% 20%;
		align-items: center;
		justify-content: center;
		width: 80%;
		max-width: 800px;
		text-wrap: wrap;
		background-color: rgba(255, 255, 255, 0.2);
		border-radius: 10px;
		box-shadow: 1px 4px 10px rgba(0, 0, 0, 0.5);
		width: auto;
		height: auto;
		min-height: 40%;
	}
  </style>
  
  <div class="search-container">
	<input type="text" bind:value={query} placeholder="Ask me anything for I am your God...">
	<button on:click={handleSearch}>Search</button>
  </div>
  {#if result}
	<div class="result-container">
		<p>{@html result}</p>
	</div>
  {:else}
	<div class="empty-result-container">
		<p>{result}</p>
	</div>
  {/if}
  