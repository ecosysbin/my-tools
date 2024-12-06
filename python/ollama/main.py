from langchain_community.llms import Ollama
llm = Ollama(model="gemma2:2b")
result = llm.invoke("Why is the sky blue?")
print(result)