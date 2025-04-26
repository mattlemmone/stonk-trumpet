import json
import logging
import openai

SYSTEM_PROMPT = """You are an AI assistant specialized in financial market sentiment analysis. 
        Analyze the following social media post (may contain HTML) regarding its potential impact on the stock market. 
        Classify the sentiment as 'positive' (number go up), 'negative' (number go down), or 'neutral' (no impact).
        Also, determine if the potential impact is 'significant' (true or false) and likely (as in, you're 99% sure, it's obvious to a savvy investor) to have immediate impact on the market.
        Additionally, provide a brief reasoning (1-2 sentences) explaining why you believe this will influence the broader market.
        
        Respond ONLY with a JSON object containing three keys: 
        - "sentiment" (string): "positive", "negative", or "neutral"
        - "significant" (boolean): true or false
        - "reasoning" (string): brief explanation of market impact
        
        Example: {"sentiment": "positive", "significant": true, "reasoning": "The announcement of tax cuts for businesses could stimulate economic growth and boost investor confidence."}"""

class SentimentAnalyzer:
    def __init__(self, api_key):
        self.api_key = api_key
        self.open_api_client = openai.OpenAI(api_key=self.api_key)
        
    def analyze(self, status_text_html):
        """
        Analyzes sentiment using OpenAI for stock market impact.
        
        Returns:
            tuple: (sentiment, significant, reasoning)
                sentiment (str): 'positive', 'negative', or 'neutral'
                significant (bool): True if the sentiment is significant, False otherwise
                reasoning (str): A brief explanation of the reasoning
        """
        log_text = status_text_html[:150] + ('...' if len(status_text_html) > 150 else '')
        logging.info(f"Analyzing HTML content: '{log_text}'")
        sentiment = "neutral"
        significant = False
        reasoning = ""


        try:
            response = self.open_api_client.chat.completions.create(
                model="gpt-4o-mini",
                messages=[
                    {"role": "system", "content": SYSTEM_PROMPT},
                    {"role": "user", "content": status_text_html}
                ],
                temperature=1.0,
                max_tokens=150,
                response_format={"type": "json_object"}
            )
            content = response.choices[0].message.content
            logging.debug(f"OpenAI raw response: {content}")
            result = json.loads(content)
            sentiment = result.get("sentiment", "neutral").lower()
            significant = result.get("significant", False)
            reasoning = result.get("reasoning", "")

            if sentiment not in ["positive", "negative", "neutral"]:
                logging.warning(f"Received unexpected sentiment value: {sentiment}. Defaulting.")
                sentiment = "neutral"
            if not isinstance(significant, bool):
                logging.warning(f"Received non-boolean significance value: {significant}. Defaulting.")
                significant = False

            logging.info(f"Analysis result: Sentiment='{sentiment}', Significant={significant}, Reasoning='{reasoning}'")

        except json.JSONDecodeError as e:
            logging.error(f"Failed to decode JSON response from OpenAI: {e} - Response: {content}")
        except (openai.APIError, openai.APITimeoutError, openai.RateLimitError, Exception) as e:
            logging.error(f"OpenAI API call failed: {e}")

        return sentiment, significant, reasoning 