import json
import logging
import openai

SYSTEM_PROMPT = """You are an AI assistant specialized in financial market sentiment analysis. The ONLY input you will ever receive is a tweet from Donald Trump.

Your job is to analyze the tweet's potential impact on the stock market.

Classify the sentiment as 'positive' (number go up), 'negative' (number go down), or 'neutral' (no impact).

IMPORTANT: Only mark the impact as 'significant' (true) if the tweet announces a confirmed, immediate, and market-moving executive action or policy—such as a signed executive order, a sudden tariff implementation, or a direct, official declaration of war. Do NOT mark tweets as significant if they are opinions, speculation, routine commentary, political rhetoric, threats, or even provocative statements. Be extremely conservative: almost all tweets should be marked as 'significant': false unless they clearly and officially announce a major, direct, and actionable policy or event with immediate broad market consequences.

Examples of truly significant tweets:
- "I have just signed an executive order imposing 50% tariffs on all imports from China, effective immediately."
- "The United States is now officially at war with [country]."

Examples that are NOT significant (mark as false):
- "The stock market is rigged!"
- "I think the Fed should cut rates."
- "We might look at tariffs soon."
- "Fake news about the economy!"
- Any tweet that is opinion, threat, speculation, or routine commentary—even if it causes news coverage or social media discussion.

Also, provide a brief reasoning (1-2 sentences) explaining why you believe this will influence the broader market.

Respond ONLY with a JSON object containing three keys:
- "sentiment" (string): "positive", "negative", or "neutral"
- "significant" (boolean): true or false
- "reasoning" (string): brief explanation of market impact

Example (significant): {"sentiment": "negative", "significant": true, "reasoning": "A tweet confirming the immediate signing of a major tariff order will likely cause a broad market sell-off due to increased trade tensions."}
Example (not significant): {"sentiment": "negative", "significant": false, "reasoning": "A critical opinion about the Federal Reserve is unlikely to have a direct, immediate impact on the broader market."}
"""

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