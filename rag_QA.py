from llama_index.embeddings.huggingface import HuggingFaceEmbedding
from llama_index.core import Settings, SimpleDirectoryReader, VectorStoreIndex
from llama_index.core.retrievers import VectorIndexRetriever
from llama_index.core.query_engine import RetrieverQueryEngine
from llama_index.core.postprocessor import SimilarityPostprocessor
import subprocess
import sys

# define settings
# import any embedding model on HF hub (https://huggingface.co/spaces/mteb/leaderboard)
Settings.embed_model = HuggingFaceEmbedding(model_name="BAAI/bge-small-en-v1.5")
# Settings.embed_model = HuggingFaceEmbedding(model_name="thenlper/gte-large") # alternative model

Settings.llm = None
Settings.chunk_size = 200
Settings.chunk_overlap = 25

# Read and Store Docs into Vector DB
# articles available here: {add GitHub repo}
documents = SimpleDirectoryReader("./Documents/QandA").load_data()
# store docs into vector DB
index = VectorStoreIndex.from_documents(documents)

# set number of docs to retreive
top_k = 1

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

# reformat response
context = "Context:\n"
for i in range(top_k):
    text_without_newlines = response.source_nodes[i].text.replace("\n", " ")  # Replace newline characters with spaces
    context = context + text_without_newlines

# onto prompting our LLM
def generate_prompt_with_context(context, query):
    prompt_template_w_context = f"""[INST]You are an AI-Dietitian.\
    Your diabetic patients asks you questions about diabetes. You end responses with the signature '– AI-Dietitian'. \

    {context}
    Please respond to the following comment briefly. Use the context above if it is helpful.

    {query}
    [/INST]
    """
    return prompt_template_w_context

prompt = generate_prompt_with_context(context, query)
print(prompt)
