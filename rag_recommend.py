from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.core import Settings, SimpleDirectoryReader, VectorStoreIndex
from llama_index.core.retrievers import VectorIndexRetriever
from llama_index.core.query_engine import RetrieverQueryEngine
from llama_index.core.postprocessor import SimilarityPostprocessor
from llama_parse import LlamaParse
from llama_index.readers.json import JSONReader
import subprocess
import sys

# define settings
# import any embedding model on HF hub (https://huggingface.co/spaces/mteb/leaderboard)
Settings.embed_model = HuggingFaceEmbedding(model_name="thenlper/gte-large")
# Settings.embed_model = HuggingFaceEmbedding(model_name="thenlper/gte-large") # alternative model

Settings.llm = None
Settings.chunk_size = 1000
Settings.chunk_overlap = 0


parser = LlamaParse(
	api_key="llx-ZdvzZ0r0W1cxC1SkXXgibSBQ8FrsQGUgKpuQKoXIxHZnxmvm",
	result_type="text",
	verbose=True,
	)

file_extractor = {".pdf": parser}
documents = SimpleDirectoryReader("./Documents/MealRec", file_extractor=file_extractor).load_data() 

# store docs into vector DB
index = VectorStoreIndex.from_documents(documents)

# set number of docs to retreive
top_k = 3

# configure retriever
retriever = VectorIndexRetriever(
    index=index,
    similarity_top_k=top_k,
)

# assemble query engine
query_engine = RetrieverQueryEngine(
    retriever=retriever,
    node_postprocessors=[SimilarityPostprocessor(similarity_cutoff=0.5)],
)

# query documents
query = sys.argv[1]
response = query_engine.query(query)

# # reformat response
# # context = "Context:\n"
# for i in range(top_k):
#     text_without_newlines = response.source_nodes[i].text.replace("\n", " ")  # Replace newline characters with spaces
#     context = context + text_without_newlines

# onto prompting our LLM
def generate_prompt_with_context(response):
    prompt_template_w_context = f"""[INST]You are an AI-Dietitian.\
    Your diabetic patients asks you for meal recommendations\

    {response}
    
    Only output the filepath of your meal recommendation.
    """
    return prompt_template_w_context

prompt = generate_prompt_with_context(response)
print(prompt)
