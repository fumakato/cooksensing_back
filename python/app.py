from flask import Flask, jsonify, request
import featureExtraction
import logging  # ロギングモジュールをインポート

app = Flask(__name__)

# ロギングの設定
logging.basicConfig(level=logging.ERROR)  # ロギングレベルを設定
logger = logging.getLogger(__name__)

# GETリクエスト用のルート

@app.route('/get2', methods=['GET'])
def get_example2():
    response = {
        'message': 'This is a GET request example.'
    }
    return jsonify(response)

@app.route('/feature_extraction', methods=['POST'])
def post_example1():
    data = request.get_json()
    url = data.get('url')
    response = featureExtraction.feature_extraction(url)
    # result = featureExtraction.feature_extraction()
    # response = {
    #     'result': result
    # }
    # return jsonify(response)
    return response

# カスタムエラーハンドラ
@app.errorhandler(Exception)
def handle_exception(e):
    logger.error(f"Unhandled Exception: {e}", exc_info=True)
    return jsonify({
        'error': 'サーバー内部で予期しないエラーが発生しました。',
        'code': 'INTERNAL_SERVER_ERROR'
    }), 500


# サーバを実行
if __name__ == '__main__':
    app.run(host="0.0.0.0", debug=True, port=5001)

