<script>
	import { onMount } from "svelte";

	let result = "";  
	let persona = 'Jesus';
	  // background-image: url('heavenly-background.jpg');
  let messages = [];
  let inputMessage = '';

  let socket;
	onMount(() => {
		connect();
	});

	async function handleSearch() {
	  if (query.trim()) {
		// sound.currentTime = 0;
		// sound.play();
		// Here, we stub out an API call
		const response = await fetch('http://api.local/v1/completion', {
		  method: 'POST',
		  headers: { 'Content-Type': 'application/json' },
		  body: JSON.stringify({ 'query': inputMessage, 'persona': persona })
		});
  
		const data = await response.json();
		console.log(data.answer); // Log the response for now
		result = data.answer.replace(/\n/g, '<br>');; 
	  }
	}



  function connect() {
    socket = new WebSocket('ws://api.local/v1/completion/stream');

    socket.onmessage = function (event) {
      messages = [...messages, event.data];
    };

    socket.onclose = function (event) {
      console.log('WebSocket closed:', event);
    };

    socket.onerror = function (error) {
      console.log('WebSocket error:', error);
    };
  }

  function sendMessage() {
    if (socket && inputMessage.trim() !== '') {
      const message = { 
		query: inputMessage,
		persona: persona 
	};

      socket.send(JSON.stringify(message));
      inputMessage = '';
    }
  }
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
	<select bind:value={persona}>
		<option value="Jesus">Jesus</option>
		<option value="Lao Tzu">Lao Tzu</option>
	  </select>
	<input type="text" bind:value={inputMessage} placeholder="Ask me anything...">
	<button on:click={sendMessage}>Search</button>
  </div>
  <div class="result-container">
	<ul>
		{#each messages as message}
			<li>{@html message}</li>
		{/each}
	</ul>
  </div>
  <!--
  {#if result}
	<div class="result-container">
		<p>{@html result}</p>
	</div>
  {:else}
	<div class="empty-result-container">
		<p>{result}</p>
	</div>
  {/if}
-->
  