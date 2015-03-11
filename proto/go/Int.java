package go;

public class Int {

	public int value;

	public Int(){ }

	public Int(int value){
		this.value = value;
	}

	public int value(){
		return value;
	}

	public void set(int value){
		this.value = value;
	}

	public Ptr_int addr(){
		return new IntWrapper(this);
	}

	public boolean equals(Object o){
		return o instanceof Int && ((Int)o).value == this.value;
	}
}
