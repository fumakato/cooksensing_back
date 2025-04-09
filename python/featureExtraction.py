import pandas as pd
import statistics
import heapq
import requests
import os
import logging
from flask import jsonify
from app import app  # メインアプリケーションからインポート

def feature_extraction(url):
    try:
        urldata = url
        # print(urldata)
        try:
            response = requests.get(url)
            response.raise_for_status()  # HTTPエラーが発生した場合は例外を送出
            temp_filename = "temp_acc.csv"
            with open(temp_filename, 'wb') as file:
                file.write(response.content)
        except requests.exceptions.RequestException as e:
            logging.error(f"Failed to download the file: {e}")
            return jsonify({
                'error': '指定されたURLからCSVファイルを取得できませんでした。',
                'code': 'INVALID_URL'
            }), 400

        # CSV読み込み
        try:
            wAcc = pd.read_csv(temp_filename, encoding='utf-8')
        except Exception as e:
            return jsonify({
                'error': 'CSVファイルの読み込みに失敗しました。',
                'code': 'CSV_READ_ERROR'
            }), 400

         # 必要な列があるかチェック
        for col in ['x', 'y', 'z', 'time']:
            if col not in wAcc.columns:
                return jsonify({
                    'error': f"CSVに必要な列 '{col}' が存在しません。",
                    'code': 'CSV_MISSING_COLUMN'
                }), 422

        # ノルムの計算
        norm = ( wAcc["x"]**2 + wAcc["y"]**2+ wAcc["z"]**2 )**0.5

        ma = 300
        mi = 0
        acc_threshold = 10
        time_threshold = 10
        for j in range(3):
            print("acc_threshold:"+str(acc_threshold))
            data = norm
            cutList = []
            cutListTime = []
            cutListData = []
            timeFlg = True
            timeDeff = -3000
            thFlg = True
            for i in range(len(data)):
                if (wAcc["time"][i] - timeDeff)/1000 > time_threshold: timeFlg = True
                if data[i] > acc_threshold and mi<=(wAcc["time"][i]-wAcc["time"][0])/1000<=ma:
                    if timeFlg and thFlg:
                        timeFlg = False
                        timeDeff = wAcc["time"][i]
                        cutList.append([wAcc["time"][i], data[i]])
                        cutListTime.append(wAcc["time"][i])
                        cutListData.append(data[i])
                    elif cutList[len(cutList)-1][1] < data[i]: 
                        cutList[len(cutList)-1] = [wAcc["time"][i], data[i]]
                        cutListTime[len(cutList)-1] = wAcc["time"][i]
                        cutListData[len(cutList)-1] = data[i]
                    thFlg = False
                else: thFlg = True


            if j==0: #はじめと終わりを検出
                if len(cutListData) < 2:
                    return jsonify({
                        'error': '十分なピークデータが検出できませんでした。',
                            'code': 'INSUFFICIENT_PEAKS'
                    }), 422
                print("len(cutListData)",len(cutListData))
                # 一番目と二番目に大きいデータを取得
                top_2 = heapq.nlargest(2, cutListData)
                # 一番目と二番目に大きいデータのインデックスを取得
                print("top_2",top_2)
                top_1_index = cutListData.index(top_2[0])
                top_2_index = cutListData.index(top_2[1])

                print("cutListTime",cutListTime[top_1_index]-wAcc["time"][0])
                print("cutListTime",cutListTime[top_2_index]-wAcc["time"][0])

                    # 終わりと初めが入れ替わる可能性がある。2秒の余裕を持たせる
                if cutListTime[top_1_index]-wAcc["time"][0] < cutListTime[top_2_index]-wAcc["time"][0]:
                    mi=(cutListTime[top_1_index]-wAcc["time"][0])/1000 + 2
                    ma=(cutListTime[top_2_index]-wAcc["time"][0])/1000 - 2
                elif cutListTime[top_1_index]-wAcc["time"][0] > cutListTime[top_2_index]-wAcc["time"][0]:
                    mi=(cutListTime[top_2_index]-wAcc["time"][0])/1000 + 2
                    ma=(cutListTime[top_1_index]-wAcc["time"][0])/1000 - 2

                time_threshold=0.2

            if j==1: #適切な閾値を取り出す
                if len(cutListData) < 10:
                    return jsonify({
                        'error': '十分なピークデータが検出できませんでした。',
                        'code': 'INSUFFICIENT_PEAKS'
                    }), 422
                # 上から10つの大きいデータを取得
                top_10 = heapq.nlargest(10, cutListData)
                print("top_10",top_10)
                avetmp = sum(top_10) / 10

                acc_threshold = avetmp * 0.9
                if acc_threshold < 10.1 : acc_threshold=10.1


        print("cut:"+str(len(cutList)))

        tempo = -1
        tempoRap = []
        tempoRapCount = 0
        for i in range(len(cutList)):
            if i == 0: 
                tempoRapCount+=1
                continue
            if tempo == -1: tempo = cutList[i][0] - cutList[i-1][0]
            else: 
                timeDeff = cutList[i][0] - cutList[i-1][0]
                if timeDeff < 1500: 
                    tempo = (tempo+timeDeff)/2
                else: 
                    tempoRap.append(tempoRapCount)
                    tempoRapCount = 0
            tempoRapCount+=1
        tempoRap.append(tempoRapCount)
        tempoRapCount = 0
        for data in tempoRap: tempoRapCount+=data
        tempoRapCount/=len(tempoRap)
        # tempo/1000 1回切るのに必要な時間
        pace = 1/tempo*1000

        

        # 'aveTempo': str(tempo/1000),
        return_dict = {
            'avePace': pace,
            'aveAcc': statistics.mean(cutListData),
            'stdev': statistics.pstdev(cutListData)
        }
        
        os.remove(temp_filename)
        return jsonify(return_dict)
    
    except Exception as e:
        app.logger.error(f"Unhandled Exception: {e}", exc_info=True)
        return jsonify({
            'error': 'サーバー内部で予期しないエラーが発生しました。',
            'code': 'INTERNAL_SERVER_ERROR'
        }), 500

        